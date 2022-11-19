export const KUBERNETES_VERSION_LATEST: number = 24;
export const KUBERNETES_VERSION_MAX: number = 25;
export const KUBERNETES_VERSION_MIN: number = 23;

interface WindowEnv {
	FF_BANNER_ENABLED?: string;
	FF_BANNER_TEXT?: string;
	FF_BANNER_COLOUR?: string;
	FF_BANNER_HEIGHT?: string;
}

declare global {
	interface Window {
		_env_?: WindowEnv;
	}
}

export const FF_BANNER_ENABLED = window._env_?.FF_BANNER_ENABLED === "true";
export const FF_BANNER_TEXT = window._env_?.FF_BANNER_TEXT || "";
export const FF_BANNER_COLOUR = window._env_?.FF_BANNER_COLOUR || "";
export const FF_BANNER_HEIGHT = Number(window._env_?.FF_BANNER_HEIGHT || "0");
