<script lang="ts">
	import { radarEnabled } from '$lib/stores/rainStore';
	import { snowEnabled, temperatureEnabled, windEnabled } from '$lib/stores/weatherStore';
	import { webcamsEnabled } from '$lib/stores/webcamStore';
	import 'maplibre-gl/dist/maplibre-gl.css';
	import { onMount } from 'svelte';
	import '../../global.css';

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'c' || event.key === 'C') {
			webcamsEnabled.update((x) => {
				return !x;
			});
		}
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
		<div class="group relative flex items-center mb-4">
			<label for="toggle-webcams" class="flex items-center w-full cursor-pointer">
				<input type="checkbox" id="toggle-webcams" class="sr-only" bind:checked={$webcamsEnabled} />
				<div class="toggle-bg bg-gray-200 border-2 border-gray-200 h-6 w-11 rounded-full"></div>
				<span class="ml-3 text-gray-900 text-sm font-medium">Webcams</span>
			</label>
			<span
				class="group-hover:opacity-100 transition-opacity bg-gray-800 text-white text-xs rounded py-1 px-2 absolute left-full ml-2 opacity-0 m-2 mx-auto z-10 whitespace-nowrap"
			>
				Shortcut: C
			</span>
		</div>

		<div class="group relative flex items-center mb-4">
			<label for="toggle-snow" class="flex items-center w-full cursor-pointer">
				<input type="checkbox" id="toggle-snow" class="sr-only" bind:checked={$snowEnabled} />
				<div class="toggle-bg bg-gray-200 border-2 border-gray-200 h-6 w-11 rounded-full"></div>
				<span class="ml-3 text-gray-900 text-sm font-medium">Snow</span>
			</label>
			<span
				class="group-hover:opacity-100 transition-opacity bg-gray-800 text-white text-xs rounded py-1 px-2 absolute left-full ml-2 opacity-0 m-2 mx-auto z-10 whitespace-nowrap"
			>
				Shortcut: S
			</span>
		</div>

		<div class="group relative flex items-center mb-4">
			<label for="toggle-temperature" class="flex items-center w-full cursor-pointer">
				<input
					type="checkbox"
					id="toggle-temperature"
					class="sr-only"
					bind:checked={$temperatureEnabled}
				/>
				<div class="toggle-bg bg-gray-200 border-2 border-gray-200 h-6 w-11 rounded-full"></div>
				<span class="ml-3 text-gray-900 text-sm font-medium">Temperature</span>
			</label>
			<span
				class="group-hover:opacity-100 transition-opacity bg-gray-800 text-white text-xs rounded py-1 px-2 absolute left-full ml-2 opacity-0 m-2 mx-auto z-10 whitespace-nowrap"
			>
				Shortcut: T
			</span>
		</div>

		<div class="group relative flex items-center mb-4">
			<label for="toggle-wind" class="flex items-center w-full cursor-pointer">
				<input type="checkbox" id="toggle-wind" class="sr-only" bind:checked={$windEnabled} />
				<div class="toggle-bg bg-gray-200 border-2 border-gray-200 h-6 w-11 rounded-full"></div>
				<span class="ml-3 text-gray-900 text-sm font-medium">Wind</span>
			</label>
			<span
				class="group-hover:opacity-100 transition-opacity bg-gray-800 text-white text-xs rounded py-1 px-2 absolute left-full ml-2 opacity-0 m-2 mx-auto z-10 whitespace-nowrap"
			>
				Shortcut: W
			</span>
		</div>

		<div class="group relative flex items-center mb-4">
			<label for="toggle-weatherradar" class="flex items-center w-full cursor-pointer">
				<input
					type="checkbox"
					id="toggle-weatherradar"
					class="sr-only"
					bind:checked={$radarEnabled}
				/>
				<div class="toggle-bg bg-gray-200 border-2 border-gray-200 h-6 w-11 rounded-full"></div>
				<span class="ml-3 text-gray-900 text-sm font-medium">Radar</span>
			</label>
			<span
				class="group-hover:opacity-100 transition-opacity bg-gray-800 text-white text-xs rounded py-1 px-2 absolute left-full ml-2 opacity-0 m-2 mx-auto z-10 whitespace-nowrap"
			>
				Shortcut: R
			</span>
		</div>
	</div>
</div>

<svelte:window on:keydown={handleKeydown} />
