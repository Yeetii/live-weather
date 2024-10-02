<script lang="ts">
	import { PUBLIC_MAPTILER_KEY } from '$env/static/public';
	import type { MapStore } from '$lib/stores';
	import { MAPSTORE_CONTEXT_KEY } from '$lib/stores';
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
			style: `https://api.maptiler.com/maps/c852a07e-70f5-49c3-aebf-ad7d488e4495/style.json?key=${PUBLIC_MAPTILER_KEY}`,
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

		mapStore?.set(map);

		map.on('error', (e: Error) => {
			console.error('Map error: ', e);
		});

		map.on('load', () => {
			map.addControl(
				new maplibregl.TerrainControl({
					source: 'terrain_rgb',
					exaggeration: 1.5
				})
			);

			fetch('https://api.rainviewer.com/public/weather-maps.json')
				.then((res) => res.json())
				.then((apiData) => {
					apiData.radar.past.forEach((frame: any) => {
						map.addLayer({
							id: `rainviewer_${frame.path}`,
							type: 'raster',
							source: {
								type: 'raster',
								tiles: [apiData.host + frame.path + '/256/{z}/{x}/{y}/2/1_1.png'],
								tileSize: 256
							},
							layout: { visibility: 'none' },
							minzoom: 0,
							maxzoom: 12
						});
					});

					let i = 0;
					const interval = setInterval(() => {
						if (i > apiData.radar.past.length - 1) {
							clearInterval(interval);
							return;
						} else {
							apiData.radar.past.forEach((frame: any, index: number) => {
								map.setLayoutProperty(
									`rainviewer_${frame.path}`,
									'visibility',
									index === i || index === i - 1 ? 'visible' : 'none'
								);
							});
							if (i - 1 >= 0) {
								const frame = apiData.radar.past[i - 1];
								let opacity = 1;
								setTimeout(() => {
									const i2 = setInterval(() => {
										if (opacity <= 0) {
											return clearInterval(i2);
										}
										map.setPaintProperty(`rainviewer_${frame.path}`, 'raster-opacity', opacity);
										opacity -= 0.1;
									}, 50);
								}, 400);
							}
							i += 1;
						}
					}, 500);
				})
				.catch(console.error);
		});
	});
</script>

<div class="map w-full h-full" data-testid="map" bind:this={mapContainer}></div>

<style>
	.map {
		height: 100vh;
		width: 100vw;
	}
</style>
