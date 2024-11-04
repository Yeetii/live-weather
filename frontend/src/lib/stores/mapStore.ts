import type { AddLayerObject, GetResourceResponse, Map as MaplibreMap, Source, SourceSpecification, StyleImageInterface, StyleImageMetadata, StyleSetterOptions } from 'maplibre-gl';
import { writable } from 'svelte/store';

export const MAPSTORE_CONTEXT_KEY = 'maplibre-map-store';

export type MapStore = ReturnType<typeof createMapStore>;

// map store for maplibre-gl object
export const createMapStore = () => {
	const { set, update, subscribe } = writable<MaplibreMap>(undefined);

	const mapInitialized = writable(false);

	const subscibeMapInitialized = mapInitialized.subscribe

	const setMap = (map: MaplibreMap) => {
		set(map);
		map.on('load', () => {
			mapInitialized.set(true);
		});
	};

	const addLayer = (layer: AddLayerObject, beforeId?: string) => {
		update((state) => {
			if (state) {
				state.addLayer(layer, beforeId);
			}
			return state;
		});
	};

	const hasLayer = (id: string) => {
		let hasLayer: boolean = false;
		const unsubscriber = subscribe((state) => {
			if (state) {
				hasLayer = state.getLayer(id) != undefined;
			}
		});
		unsubscriber()
		return hasLayer
	}

	const addSource = (id: string, source: SourceSpecification) => {
		update((state) => {
			if (state) {
				state.addSource(id, source);
			}
			return state;
		});
	}

	const getSource = (id: string) => {
		let source: Source | undefined;

		const unsubscriber = subscribe((state) => {
			if (state) {
				source = state.getSource(id);
			}
		});
		unsubscriber()
		return source
	}

	const removeSource = (id: string) => {
		update((state) => {
			if (state) {
				state.removeSource(id);
			}
			return state;
		});
	}

	/**
	 * Update Maplibre's PaintProperty
	 *
	 * Note.
	 * setPaintProperty does render map canvas with new given property value.
	 * But in some cases, it does not actually update style.json object in Map instance.
	 * Because of this problem of Maplibre, the function is going to update style.json directly and call `setStyle` function.
	 *
	 * @param layerId The ID of the layer to set the paint property in.
	 * @param name The name of the paint property to set.
	 * @param value The value of the paint property to set. Must be of a type appropriate for the property, as defined in the MapLibre Style Specification.
	 * @param options Options object.
	 */
	const setPaintProperty = (
		layerId: string,
		name: string,
		value: unknown,
		options?: StyleSetterOptions
	) => {
		update((state) => {
			if (state) {
				state.setPaintProperty(layerId, name, value, options);

				const style = state.getStyle();
				const layer = style?.layers?.find((l) => l.id === layerId);
				if (layer) {
					if (!layer.paint) {
						layer.paint = {};
					}
					if (value) {
						// eslint-disable-next-line @typescript-eslint/ban-ts-comment
						// @ts-ignore
						layer.paint[name] = value;
					} else {
						// eslint-disable-next-line @typescript-eslint/ban-ts-comment
						// @ts-ignore
						if (layer.paint[name]) {
							// eslint-disable-next-line @typescript-eslint/ban-ts-comment
							// @ts-ignore
							delete layer.paint[name];
						}
					}
					state.setStyle(style);
				}
			}

			return state;
		});
	};

	/**
	 * Update Maplibre's LayoutProperty
	 *
	 * Note.
	 * setLayoutProperty does render map canvas with new given property value.
	 * But in some cases, it does not actually update style.json object in Map instance.
	 * Because of this problem of Maplibre, the function is going to update style.json directly and call `setStyle` function.
	 *
	 * @param layerId The ID of the layer to set the paint property in.
	 * @param name The name of the paint property to set.
	 * @param value The value of the paint property to set. Must be of a type appropriate for the property, as defined in the MapLibre Style Specification.
	 * @param options Options object.
	 */
	const setLayoutProperty = (
		layerId: string,
		name: string,
		value: unknown,
		options?: StyleSetterOptions
	) => {
		update((state) => {
			if (state) {
				state.setLayoutProperty(layerId, name, value, options);

				const style = state.getStyle();
				const layer = style?.layers?.find((l) => l.id === layerId);
				if (layer) {
					if (!layer.layout) {
						layer.layout = {};
					}
					if (value) {
						// eslint-disable-next-line @typescript-eslint/ban-ts-comment
						// @ts-ignore
						layer.layout[name] = value;
					} else {
						// eslint-disable-next-line @typescript-eslint/ban-ts-comment
						// @ts-ignore
						if (layer.layout[name]) {
							// eslint-disable-next-line @typescript-eslint/ban-ts-comment
							// @ts-ignore
							delete layer.layout[name];
						}
					}
					state.setStyle(style);
				}
			}

			return state;
		});
	};

	const loadImage = (url: string) => {
		let promise: Promise<GetResourceResponse<HTMLImageElement | ImageBitmap>> | null = null;
		update((state) => {
		  if (state) {
			promise = state.loadImage(url);
		  } 
		  return state;
		});
		return promise!; // Note the non-null assertion operator (!) here
	  };

	  const hasImage = (id: string) => {
		let hasImage: boolean = false;
		const unsubscriber = subscribe((state) => {
		  if (state) {
			hasImage = state.hasImage(id);
		  }
		});
		unsubscriber()
		return hasImage
	  };

	  const addImage = (id: string, image: HTMLImageElement | ImageBitmap | ImageData | {
		width: number;
		height: number;
		data: Uint8Array | Uint8ClampedArray;
		} | StyleImageInterface, options?: Partial<StyleImageMetadata>) => {
		update((state) => {
		  if (state) {
			state.addImage(id, image, options);
		  }
		  return state;
		});
	  };

	return {
		subscribe,
		update,
		setMap,
		setPaintProperty,
		setLayoutProperty,
		addLayer,
		hasLayer,
		addSource,
		getSource,
		removeSource,
		subscibeMapInitialized,
		loadImage,
		hasImage,
		addImage
	};
};
