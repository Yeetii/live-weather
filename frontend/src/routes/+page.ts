import type { PageLoad } from './$types';

export const load: PageLoad = async () => {
	const title = 'Årevädret';
	const site_name = 'Årevädret';
	const site_description = 'Årevädret - Live!';

	return {
		title,
		site_name,
		site_description
	};
};
