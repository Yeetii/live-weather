<script lang="ts">
	import type { MapStore } from '$lib/stores';
	import { MAPSTORE_CONTEXT_KEY } from '$lib/stores';
	import { radarData, type CloudData } from '$lib/stores/rainStore';
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

	const animateWeather = (data: CloudData) => {
		let i = 0;
		const interval = setInterval(() => {
			if (i > data.radar.length - 1) {
				i = 0;
				data.radar.forEach((frame) => {
					mapStore.setPaintProperty(`rainviewer_${frame.path}`, 'raster-opacity', 1);
				});
				return;
			} else {
				data.radar.forEach((frame, index: number) => {
					mapStore.setLayoutProperty(
						`rainviewer_${frame.path}`,
						'visibility',
						index === i || index === i - 1 ? 'visible' : 'none'
					);
				});
				if (i - 1 >= 0) {
					const frame = data.radar[i - 1];
					let opacity = 1;
					setTimeout(() => {
						const i2 = setInterval(() => {
							if (opacity <= 0) {
								return clearInterval(i2);
							}
							mapStore.setPaintProperty(`rainviewer_${frame.path}`, 'raster-opacity', opacity);
							opacity -= 0.1;
						}, 20);
					}, 150);
				}
				i += 1;
			}
		}, 300);
	};

	type WebcamData = {
		imageUrl: string;
		coordinates: number[];
	};

	function addWebcamsToMap(map: maplibregl.Map, webcams: WebcamData[]) {
		webcams.forEach((webcam, index) => {
			const imageId = `webcam-image-${index}`; // Unique image ID for each image

			map
				.loadImage(webcam.imageUrl)
				.then((image) => {
					// Add each image with a unique ID
					if (!map.hasImage(imageId)) {
						map.addImage(imageId, image.data);
					}

					// Add GeoJSON source for each webcam
					const sourceId = `webcam-point-${index}`;
					if (!map.getSource(sourceId)) {
						map.addSource(sourceId, {
							type: 'geojson',
							data: {
								type: 'FeatureCollection',
								features: [
									{
										type: 'Feature',
										geometry: {
											type: 'Point',
											coordinates: webcam.coordinates
										},
										properties: {
											name: `Webcam ${index + 1}`
										}
									}
								]
							}
						});
					}

					// Add a layer to display the image at each webcam's coordinates
					const layerId = `webcam-layer-${index}`;
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
					console.error(`Error loading image for webcam ${index + 1}:`, error);
				});
		});
	}

	onMount(() => {
		radarData.subscribe((data) => {
			if (!data) return;

			data.radar.forEach((frame) => {
				mapStore.addLayer({
					id: `rainviewer_${frame.path}`,
					type: 'raster',
					source: {
						type: 'raster',
						tiles: [data.host + frame.path + '/256/{z}/{x}/{y}/2/1_1.png'],
						tileSize: 256
					},
					layout: { visibility: 'none' },
					minzoom: 0,
					maxzoom: 12
				});
			});
			animateWeather(data);
		});
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
			style: `https://api.maptiler.com/maps/c852a07e-70f5-49c3-aebf-ad7d488e4495/style.json?key=KxXGPUn8leqAeKO3GqWn`,
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
			map.on('zoom', (e: any) => {
				console.log(e.target.getZoom());
			});

			map.addControl(
				new maplibregl.TerrainControl({
					source: 'terrain_rgb',
					exaggeration: 1.5
				})
			);

			const webcams = [
				{
					imageUrl:
						'https://api.trafikinfo.trafikverket.se/v2/Images/RoadConditionCamera_39635528.Jpeg?type=fullsize&maxage=140',
					coordinates: [13.206262037974694, 63.25443937020817]
				},
				{
					imageUrl:
						'https://api.trafikinfo.trafikverket.se/v2/Images/RoadConditionCamera_39635384.Jpeg?type=fullsize&maxage=140',
					coordinates: [12.702011476723557, 63.36705212550013]
				},
				{
					imageUrl:
						'https://api.trafikinfo.trafikverket.se/v2/Images/RoadConditionCamera_39635520.Jpeg?type=fullsize&maxage=140',
					coordinates: [12.702011476723557, 63.36705212550013]
				},
				{
					imageUrl:
						'https://api.trafikinfo.trafikverket.se/v2/Images/RoadConditionCamera_39626819.Jpeg?type=fullsize&maxage=140',
					coordinates: [12.407996503920517, 63.519546112666376]
				}
			];
			addWebcamsToMap(map, webcams);
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
