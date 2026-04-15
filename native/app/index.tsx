import { Redirect } from 'expo-router'
import { useEffect, useState } from 'react'
import { storage } from '../lib/storage'

export default function Index() {
  const [isLoading, setIsLoading] = useState(true)
  const [hasToken, setHasToken] = useState(false)

  useEffect(() => {
    storage.getToken().then((token) => {
      setHasToken(!!token)
      setIsLoading(false)
    })
  }, [])

  if (isLoading) return null

  return <Redirect href={hasToken ? '/(app)/todos' : '/(auth)/login'} />
}
