import { dev } from "$app/environment";
import type { Feature, GeoJsonProperties, Geometry } from "geojson";
import { writable } from "svelte/store";

export const webcamStore = writable<Webcam[]>([]);
export const webcamsEnabled = writable(true);

export interface Webcam {
  url: string;
  location: Feature<Geometry, GeoJsonProperties>;
}

const url = dev
			? 'http://localhost:8080/fetchWebcams'
			: 'https://api.weather.erikmagnusson.com/fetchWebcams';

fetch(url)
  .then((response) => response.json())
  .then((data: Webcam[]) => {
    webcamStore.set(data)
  });