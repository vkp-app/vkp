export const KUBERNETES_VERSION_LATEST: number = 26;
export const KUBERNETES_VERSION_MAX: number = 27;
export const KUBERNETES_VERSION_MIN: number = 23;

interface WindowEnv {
	FF_BANNER_ENABLED?: string;
	FF_BANNER_TEXT?: string;
	FF_BANNER_COLOUR?: string;
	FF_BANNER_HEIGHT?: string;
	FF_ALLOW_HA?: string;
	FF_HELP_URL?: string;
}

declare global {
	interface Window {
		_env_?: WindowEnv;
	}
}

export const FF_BANNER_ENABLED = window._env_?.FF_BANNER_ENABLED === "true";
export const FF_BANNER_TEXT = window._env_?.FF_BANNER_TEXT || "FF_BANNER_TEXT IS NOT SET!";
export const FF_BANNER_COLOUR = window._env_?.FF_BANNER_COLOUR || "#4CAF50";
export const FF_BANNER_HEIGHT = Number(window._env_?.FF_BANNER_HEIGHT || "32") || 32;

export const FF_ALLOW_HA = window._env_?.FF_ALLOW_HA === "true";

export const FF_HELP_URL = window._env_?.FF_HELP_URL || "https://vkp-app.github.io/docs/";
