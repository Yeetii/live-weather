import type { PageLoad } from './$types';

export const load: PageLoad = async () => {
	const title = 'Peakshunters';
	const site_name = 'Peakshunters';
	const site_description = 'Visualizes peaks summited in Strava activities';

	return {
		title,
		site_name,
		site_description
	};
};
