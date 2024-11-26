<script lang="ts">
	import type { MapStore } from '$lib/stores';
	import { MAPSTORE_CONTEXT_KEY } from '$lib/stores';
	import { radarData, radarEnabled } from '$lib/stores/rainStore';
	import 'maplibre-gl/dist/maplibre-gl.css';
	import { getContext, onDestroy, onMount } from 'svelte';
	import '../../global.css';

	let mapStore: MapStore = getContext(MAPSTORE_CONTEXT_KEY);

	type RadarLayer = {
		id: string;
		timestamp: number;
	};

	var radarLayers: RadarLayer[];
	var interval: NodeJS.Timeout;
	var currentTimestamp: number;
	var showTimestamp: boolean;

	radarEnabled.subscribe((status) => {
		if (!radarLayers) {
			radarEnabled.set(false);
			return;
		}

		showTimestamp = status;
		if (status) {
			animateWeather();
		} else {
			disableRadar();
		}
	});

	const disableRadar = () => {
		if (radarLayers) {
			radarLayers.forEach((layer) => {
				mapStore.setLayoutProperty(layer.id, 'visibility', 'none');
			});
		}
		clearInterval(interval);
	};

	const animateWeather = () => {
		let i = 0;

		interval = setInterval(() => {
			currentTimestamp = radarLayers[i].timestamp;
			mapStore.setLayoutProperty(
				radarLayers[(i + radarLayers.length - 1) % radarLayers.length].id,
				'visibility',
				'none'
			);
			mapStore.setLayoutProperty(radarLayers[i].id, 'visibility', 'visible');
			i = (i + 1) % radarLayers.length;
		}, 600);
	};

	const timestampToString = (timestamp: number) => {
		if (!timestamp) return '';

		const date = new Date(timestamp * 1000);
		const hours = date.getHours().toString().padStart(2, '0');
		const minutes = date.getMinutes().toString().padStart(2, '0');
		return `${hours}:${minutes}`;
	};

	onMount(() => {
		radarData.subscribe((data) => {
			if (!data) return;

			mapStore.subscribe((map) => {
				if (!map) return;
			});

			mapStore.subscibeMapInitialized((value) => {
				if (!value) return;
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
				radarLayers = data.radar.map((frame) => ({
					id: `rainviewer_${frame.path}`,
					timestamp: frame.time
				}));
			});
		});
	});

	onDestroy(() => {
		disableRadar();
	});
</script>

{#if showTimestamp}
	<div class="absolute bottom-5 left-5 z-50 font-bold text-lg">
		{timestampToString(currentTimestamp)}
	</div>
{/if}
