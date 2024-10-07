<script lang="ts">
	import type { MapStore } from '$lib/stores';
	import { MAPSTORE_CONTEXT_KEY } from '$lib/stores';
	import { radarData } from '$lib/stores/rainStore';
	import 'maplibre-gl/dist/maplibre-gl.css';
	import { getContext, onMount } from 'svelte';
	import '../../global.css';

	let mapStore: MapStore = getContext(MAPSTORE_CONTEXT_KEY);

	var radarLayerIds: string[];
	var interval: NodeJS.Timeout;
	var radarEnabled = false;

	const toggleRadar = (_e: Event) => {
		if (!radarLayerIds) {
			radarEnabled = false;
			return;
		}

		if (!radarEnabled) {
			animateWeather();
		} else {
			disableRadar();
		}
	};

	const disableRadar = () => {
		radarLayerIds.forEach((id) => {
			mapStore.setLayoutProperty(id, 'visibility', 'none');
		});
		clearInterval(interval);
	};

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'w' || event.key === 'W') {
			toggleRadar(event);
			radarEnabled = !radarEnabled;
		}
	}

	const animateWeather = () => {
		let i = 0;

		interval = setInterval(() => {
			radarLayerIds.forEach((layerId, index: number) => {
				mapStore.setLayoutProperty(layerId, 'visibility', index === i ? 'visible' : 'none');
			});
		}, 600);
	};

	onMount(() => {
		radarData.subscribe((data) => {
			if (!data) return;

			mapStore.subscribe((map) => {
				if (!map) return;
			});

			mapStore.subscibeMapInitialized((value) => {
				if (!value) return;
				console.log('jsjsj');
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
				radarLayerIds = data.radar.map((frame) => `rainviewer_${frame.path}`);
			});
		});
	});
</script>

<div class="absolute top-5 left-5 z-50">
	<div class="max-w-lg mx-auto">
		<label for="toggle-weatherradar" class="flex items-center cursor-pointer relative mb-4">
			<input
				type="checkbox"
				id="toggle-weatherradar"
				class="sr-only"
				bind:checked={radarEnabled}
				on:change={toggleRadar}
			/>
			<div class="toggle-bg bg-gray-200 border-2 border-gray-200 h-6 w-11 rounded-full"></div>
			<span class="ml-3 text-gray-900 text-sm font-medium">Weather radar</span>
		</label>
	</div>
</div>

<svelte:window on:keydown={handleKeydown} />
