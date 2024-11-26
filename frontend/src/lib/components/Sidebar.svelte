<script lang="ts">
	import { radarEnabled } from '$lib/stores/rainStore';
	import 'maplibre-gl/dist/maplibre-gl.css';
	import { onMount } from 'svelte';
	import '../../global.css';

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'w' || event.key === 'W') {
			radarEnabled.update((x) => {
				return !x;
			});
		}
	}

	onMount(() => {});
</script>

<div class="absolute top-5 left-5 z-50">
	<div class="max-w-lg mx-auto">
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
