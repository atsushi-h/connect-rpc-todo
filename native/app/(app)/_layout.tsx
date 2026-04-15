import { router, Stack } from 'expo-router'
import { useEffect, useState } from 'react'
import { storage } from '../../lib/storage'

export default function AppLayout() {
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    storage.getToken().then((token) => {
      if (!token) {
        router.replace('/(auth)/login')
      }
      setIsLoading(false)
    })
  }, [])

  if (isLoading) return null

  return <Stack />
}
