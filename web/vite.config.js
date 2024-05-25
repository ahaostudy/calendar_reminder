import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

// arco.design
import { vitePluginForArco } from '@arco-plugins/vite-vue'

// https://vitejs.dev/config/
export default ({ mode }) => {
  const env = { ...loadEnv(mode, process.cwd()) }

  return defineConfig({
    plugins: [
      vue(),
      vitePluginForArco({
        style: 'css'
      })
    ],
    resolve: {
      alias: {
        '@': resolve(__dirname, 'src'),
        '@c': resolve(__dirname, 'src/components')
      }
    },
  })
}
