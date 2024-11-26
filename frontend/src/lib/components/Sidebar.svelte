<script lang="ts">
	import { radarEnabled } from '$lib/stores/rainStore';
	import { snowEnabled, temperatureEnabled, windEnabled } from '$lib/stores/weatherStore';
	import 'maplibre-gl/dist/maplibre-gl.css';
	import { onMount } from 'svelte';
	import '../../global.css';

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'r' || event.key === 'R') {
			radarEnabled.update((x) => {
				return !x;
			});
		}
		if (event.key === 's' || event.key === 'S') {
			snowEnabled.update((x) => {
				return !x;
			});
		}
		if (event.key === 't' || event.key === 'T') {
			temperatureEnabled.update((x) => {
				return !x;
			});
		}
		if (event.key === 'w' || event.key === 'W') {
			windEnabled.update((x) => {
				return !x;
			});
		}
	}

	onMount(() => {});
</script>

<div class="absolute top-5 left-5 z-50">
	<div class="max-w-lg mx-auto">
		<label for="toggle-snow" class="flex items-center cursor-pointer relative mb-4">
			<input type="checkbox" id="toggle-snow" class="sr-only" bind:checked={$snowEnabled} />
			<div class="toggle-bg bg-gray-200 border-2 border-gray-200 h-6 w-11 rounded-full"></div>
			<span class="ml-3 text-gray-900 text-sm font-medium">Snow data</span>
		</label>
		<label for="toggle-temperature" class="flex items-center cursor-pointer relative mb-4">
			<input
				type="checkbox"
				id="toggle-temperature"
				class="sr-only"
				bind:checked={$temperatureEnabled}
			/>
			<div class="toggle-bg bg-gray-200 border-2 border-gray-200 h-6 w-11 rounded-full"></div>
			<span class="ml-3 text-gray-900 text-sm font-medium">Temperature data</span>
		</label>
		<label for="toggle-wind" class="flex items-center cursor-pointer relative mb-4">
			<input type="checkbox" id="toggle-wind" class="sr-only" bind:checked={$windEnabled} />
			<div class="toggle-bg bg-gray-200 border-2 border-gray-200 h-6 w-11 rounded-full"></div>
			<span class="ml-3 text-gray-900 text-sm font-medium">Wind data</span>
		</label>
		<label for="toggle-weatherradar" class="flex items-center cursor-pointer relative mb-4">
			<input
				type="checkbox"
				id="toggle-weatherradar"
				class="sr-only"
				bind:checked={$radarEnabled}
			/>
			<div class="toggle-bg bg-gray-200 border-2 border-gray-200 h-6 w-11 rounded-full"></div>
			<span class="ml-3 text-gray-900 text-sm font-medium">Weather radar</span>
		</label>
	</div>
</div>

<svelte:window on:keydown={handleKeydown} />
