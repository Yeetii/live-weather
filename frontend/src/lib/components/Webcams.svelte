<script lang="ts">
	import { MAPSTORE_CONTEXT_KEY, type MapStore } from '$lib/stores';
	import { webcamsEnabled, webcamStore, type Webcam } from '$lib/stores/webcamStore';
	import { getContext } from 'svelte';

	let mapStore: MapStore = getContext(MAPSTORE_CONTEXT_KEY);
	let layerIds: string[] = [];

	mapStore.subscibeMapInitialized((initialized) => {
		if (!initialized) return;

		webcamStore.subscribe((webcams) => {
			addWebcamsToMap(webcams).then(() => {
				webcamsEnabled.subscribe((status) => {
					if (status) {
						layerIds.forEach((id) => {
							mapStore.setLayoutProperty(id, 'visibility', 'visible');
						});
					} else {
						layerIds.forEach((id) => {
							mapStore.setLayoutProperty(id, 'visibility', 'none');
						});
					}
				});
			});
		});
	});

	async function addWebcamsToMap(data: Webcam[]): Promise<void> {
		await Promise.all(
			data.map(async (webcam) => {
				const imageId = webcam.location.id as string;

				try {
					const image = await mapStore.loadImage(webcam.url);
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
						layerIds.push(layerId);
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
								],
								visibility: 'none'
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
				} catch (error) {
					console.error(`Error loading image for webcam ${imageId}:`, error);
				}
			})
		);
	}
</script>
