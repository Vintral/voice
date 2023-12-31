import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import './globals.css'

import Footer from "@/app/components/footer"
import Header from "@/app/components/header"

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: 'Create Next App',
  description: 'Generated by create next app',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  // return (
  //   <html lang="en">
  //     <body className={inter.className}>{children}</body>
  //   </html>
  // )

  return (
    <html lang="en">
      <body className={inter.className}>
        <div className="flex flex-col justify-between w-full h-full min-h-screen">
          <Header />          
          {children}
          <Footer />
        </div>
      </body>
    </html>
  )  
}
