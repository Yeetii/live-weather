import { initializeApp } from "firebase/app";
import { collection, getFirestore, onSnapshot } from "firebase/firestore";
import type { Feature } from "geojson";
import { writable } from "svelte/store";

export const weatherStore = writable<Feature[]>([]);
export const snowEnabled = writable(true);
export const temperatureEnabled = writable(false);
export const windEnabled = writable(false);

const app = initializeApp({projectId: "live-weather-eefc5"});
const db = getFirestore(app);

function subscribeToCollection() {
  const collectionRef = collection(db, "weatherObservations");

  const unsubscribe = onSnapshot(collectionRef, (snapshot) => {
    const documents = snapshot.docs.map((doc) => ({
      id: doc.id,
      ...doc.data(),
    }));
    weatherStore.set(documents as Feature[]);
  });
  return unsubscribe;
}

subscribeToCollection();