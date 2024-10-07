<script lang="ts">
	import { dev } from '$app/environment';
	import type { MapStore } from '$lib/stores';
	import { MAPSTORE_CONTEXT_KEY } from '$lib/stores';
	import type { Feature, GeoJsonProperties, Geometry } from 'geojson';
	import maplibregl, {
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
	import WeatherRadar from './WeatherRadar.svelte';
	import WeatherMeasurements from './WeatherMeasurements.svelte';
	
	let mapStore: MapStore = getContext(MAPSTORE_CONTEXT_KEY);

	let mapContainer: HTMLDivElement;

	interface Webcam {
		url: string;
		location: Feature<Geometry, GeoJsonProperties>;
	}

	function addWebcamsToMap(map: maplibregl.Map) {
		const url = dev
			? 'http://localhost:8080/fetchWebcams'
			: 'https://api.weather.erikmagnusson.com/fetchWebcams';
		fetch(url)
			.then((response) => response.json())
			.then((data: Webcam[]) => {
				data.forEach((webcam) => {
					const imageId = webcam.location.id as string;

					map
						.loadImage(webcam.url)
						.then((image) => {
							if (!map.hasImage(imageId)) {
								map.addImage(imageId, image.data);
							}

							const sourceId = `webcam-point-${imageId}`;
							if (!map.getSource(sourceId)) {
								map.addSource(sourceId, {
									type: 'geojson',
									data: {
										type: 'FeatureCollection',
										features: [webcam.location]
									}
								});
							}

							const layerId = `webcam-layer-${imageId}`;
							if (!map.getLayer(layerId)) {
								map.addLayer({
									id: layerId,
									type: 'symbol',
									source: sourceId,
									layout: {
										'icon-image': imageId,
										'icon-size': [
											'interpolate',
											['linear'],
											['zoom'],
											0,
											0, // At zoom level 0, size is 0
											10,
											0.01, // At zoom level 10, size is 0.01
											22,
											1 // At zoom level 22, size is 1
										]
									}
								});
							}
						})
						.catch((error: any) => {
							console.error(`Error loading image for webcam ${imageId}:`, error);
						});
				});
			});
	}

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

		map.on('load', () => {
			addWebcamsToMap(map);

			
		});
	});
</script>

<div class="map w-full h-full" data-testid="map" bind:this={mapContainer}>
	<WeatherRadar />
	<WeatherMeasurements />
</div>

<style>
	.map {
		height: 100vh;
		width: 100vw;
	}
</style>
