import '../styles/globals.css'
import type {AppProps} from 'next/app'
import localFont from '@next/font/local'

const notoSansJP = localFont({
  variable: '--font-noto-sans-jp',
  src: './NotoSansJP-VF.woff2'
})

export default function App({Component, pageProps}: AppProps) {
  return (
      <div data-theme='ictsc' className={`${notoSansJP.className}`}>
        <Component {...pageProps} />
      </div>
  )
}
