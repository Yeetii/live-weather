import { initializeApp } from "firebase/app";
import { getFirestore, collection, onSnapshot } from "firebase/firestore";
import type { Feature } from "geojson";
import { writable } from "svelte/store";

export const weatherStore = writable<Feature[]>([]);

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