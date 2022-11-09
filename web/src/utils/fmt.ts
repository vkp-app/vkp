const k = 1024;
const sizes = ["Bytes", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"];

/**
 * Adapter from https://stackoverflow.com/a/18650828
 */
export const formatBytes = (bytes: number, decimals = 2): string => {
	if (!+bytes || !Number.isFinite(bytes))
		return "0 Bytes";

	const dm = decimals < 0 ? 0 : decimals;

	const i = Math.floor(Math.log(bytes) / Math.log(k));
	return `${parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`;
}