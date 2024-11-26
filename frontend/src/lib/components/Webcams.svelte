<script lang="ts">
	import { dev } from '$app/environment';
	import { MAPSTORE_CONTEXT_KEY, type MapStore } from '$lib/stores';
	import type { Feature, GeoJsonProperties, Geometry } from 'geojson';
	import { getContext } from 'svelte';

	let mapStore: MapStore = getContext(MAPSTORE_CONTEXT_KEY);

	interface Webcam {
		url: string;
		location: Feature<Geometry, GeoJsonProperties>;
	}

	mapStore.subscibeMapInitialized((initialized) => {
		if (!initialized) return;

		addWebcamsToMap();
	});

	function addWebcamsToMap() {
		const url = dev
			? 'http://localhost:8080/fetchWebcams'
			: 'https://api.weather.erikmagnusson.com/fetchWebcams';
		fetch(url)
			.then((response) => response.json())
			.then((data: Webcam[]) => {
				data.forEach((webcam) => {
					const imageId = webcam.location.id as string;

					mapStore
						.loadImage(webcam.url)
						.then((image) => {
							if (!mapStore.hasImage(imageId)) {
								mapStore.addImage(imageId, image.data);
							}

							const sourceId = `webcam-point-${imageId}`;
							if (!mapStore.getSource(sourceId)) {
								mapStore.addSource(sourceId, {
									type: 'geojson',
									data: {
										type: 'FeatureCollection',
										features: [webcam.location]
									}
								});
							}

							const layerId = `webcam-layer-${imageId}`;
							if (!mapStore.hasLayer(layerId)) {
								mapStore.addLayer({
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
								mapStore.update((map) => {
									map.on('click', layerId, (e) => {
										map.flyTo({
											// @ts-ignore
											center: e.features[0].geometry.coordinates,
											zoom: 14
										});
									});
									map.on('mouseenter', layerId, () => {
										map.getCanvas().style.cursor = 'pointer';
									});

									map.on('mouseleave', layerId, () => {
										map.getCanvas().style.cursor = '';
									});
									return map;
								});
							}
						})
						.catch((error: any) => {
							console.error(`Error loading image for webcam ${imageId}:`, error);
						});
				});
			});
	}
</script>
