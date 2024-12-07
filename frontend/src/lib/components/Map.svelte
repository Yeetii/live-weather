<script lang="ts">
	import type { MapStore } from '$lib/stores';
	import { MAPSTORE_CONTEXT_KEY } from '$lib/stores';
	import {
		AttributionControl,
		GeolocateControl,
		LngLat,
		Map,
		NavigationControl,
		ScaleControl
	} from 'maplibre-gl';
	import 'maplibre-gl/dist/maplibre-gl.css';
	import { getContext, onMount } from 'svelte';
	import '../../global.css';
	import Sidebar from './Sidebar.svelte';
	import WeatherMeasurements from './WeatherMeasurements.svelte';
	import WeatherRadar from './WeatherRadar.svelte';
	import Webcams from './Webcams.svelte';

	let mapStore: MapStore = getContext(MAPSTORE_CONTEXT_KEY);

	let mapContainer: HTMLDivElement;

	onMount(() => {
		var center = new LngLat(13.0509, 63.41698);
		var zoom = 12;
		const mapCenter = localStorage.getItem('mapCenter');
		const mapZoom = localStorage.getItem('mapZoom');
		if (mapCenter && mapZoom) {
			center = JSON.parse(mapCenter);
			zoom = parseFloat(mapZoom);
		}

		const map = new Map({
			container: mapContainer,
			style: `https://api.maptiler.com/maps/basic-v2/style.json?key=KxXGPUn8leqAeKO3GqWn`,
			center: center,
			zoom: zoom,
			hash: true,
			attributionControl: false
		});
		map.addControl(new NavigationControl({}), 'top-right');
		map.addControl(
			new GeolocateControl({
				positionOptions: { enableHighAccuracy: true },
				trackUserLocation: true
			}),
			'top-right'
		);
		map.addControl(new ScaleControl({ maxWidth: 80, unit: 'metric' }), 'bottom-right');
		map.addControl(new AttributionControl({ compact: true }), 'bottom-right');

		mapStore?.setMap(map);

		map.on('error', (e: Error) => {
			console.error('Map error: ', e);
		});
	});
</script>

<div class="map w-full h-full" data-testid="map" bind:this={mapContainer}>
	<Sidebar />
	<WeatherRadar />
	<WeatherMeasurements />
	<Webcams />
</div>

<style>
	.map {
		height: 100vh;
		width: 100vw;
	}
</style>
