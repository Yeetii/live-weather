import { writable } from "svelte/store"

interface Rainviewer {
    version: string
    generated: number
    host: string
    radar: Radar
    satellite: Satellite
  }
  
interface Radar {
    past: TimePath[]
    nowcast: TimePath[]
  }
  
  export interface TimePath {
    time: number
    path: string
  }
  
  export interface Satellite {
    infrared: TimePath[]
  }

  export interface CloudData {
    host: string
    radar: TimePath[]
    infrared: TimePath[]
  }

  export const radarData = writable<CloudData>();
  export const radarEnabled = writable(false);

fetch('https://api.rainviewer.com/public/weather-maps.json')
    .then((res) => res.json())
    .then((apiData: Rainviewer) => {
        radarData.set(
            {host: apiData.host, radar: [...apiData.radar.past, ...apiData.radar.nowcast], infrared: apiData.satellite.infrared}
        )
    })
    .catch(console.error);