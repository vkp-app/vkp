import {resolve} from "path";
import {defineConfig} from "vite";
import react from "@vitejs/plugin-react-swc";

// https://vitejs.dev/config/
export default defineConfig({
	plugins: [react()],
	build: {
		rollupOptions: {
			input: {
				login: resolve(__dirname, "src/templates/index.html")
			}
		}
	}
})
