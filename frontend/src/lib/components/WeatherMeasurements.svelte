<script lang="ts">
	import { MAPSTORE_CONTEXT_KEY, type MapStore } from '$lib/stores';
	import {
		snowEnabled,
		temperatureEnabled,
		weatherStore,
		windEnabled
	} from '$lib/stores/weatherStore';
	import type { GeoJSONSource } from 'maplibre-gl';
	import { getContext } from 'svelte';

	let mapStore: MapStore = getContext(MAPSTORE_CONTEXT_KEY);

	mapStore.subscibeMapInitialized((initialized) => {
		if (!initialized) return;

		weatherStore.subscribe((data) => {
			if (!data) return;

			var source = mapStore.getSource('weatherMeasurement-source') as GeoJSONSource;
			if (source) {
				source.setData({
					type: 'FeatureCollection',
					features: data
				});
				return;
			}

			mapStore.addSource('weatherMeasurement-source', {
				type: 'geojson',
				data: {
					type: 'FeatureCollection',
					features: data
				}
			});

			addMeasurement('windSpeed_ms', 'm/s');
			addMeasurement('temperature_c', '°C');
			addMeasurement('snowDepth_cm', 'cm');

			snowEnabled.subscribe((status) => {
				if (status) {
					mapStore.setLayoutProperty('snowDepth_cm', 'visibility', 'visible');
				} else {
					mapStore.setLayoutProperty('snowDepth_cm', 'visibility', 'none');
				}
			});

			windEnabled.subscribe((status) => {
				if (status) {
					mapStore.setLayoutProperty('windSpeed_ms', 'visibility', 'visible');
				} else {
					mapStore.setLayoutProperty('windSpeed_ms', 'visibility', 'none');
				}
			});

			temperatureEnabled.subscribe((status) => {
				if (status) {
					mapStore.setLayoutProperty('temperature_c', 'visibility', 'visible');
				} else {
					mapStore.setLayoutProperty('temperature_c', 'visibility', 'none');
				}
			});
		});
	});

	function addMeasurement(attribute: string, unit: string) {
		mapStore.addLayer({
			id: attribute,
			type: 'symbol',
			source: 'weatherMeasurement-source',
			layout: {
				'text-field': [
					'concat',
					['get', attribute],
					['case', ['!=', ['to-string', ['get', attribute]], ''], unit, '']
				],
				'text-size': 18,
				'text-anchor': 'top',
				'text-justify': 'auto',
				'text-variable-anchor-offset': [
					'top',
					[0, -2],
					'left',
					[-2, 0],
					'bottom',
					[0, 2],
					'right',
					[2, 0]
				],
				'text-radial-offset': ['interpolate', ['linear'], ['zoom'], 5, 0, 10, 2, 15, 3],
				visibility: 'none'
			},
			paint: {
				'text-color': '#000'
			}
		});
	}
</script>
