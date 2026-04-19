import { createEnv } from '@t3-oss/env-core'
import { z } from 'zod'

export const env = createEnv({
  clientPrefix: 'EXPO_PUBLIC_',
  server: {},
  client: {
    EXPO_PUBLIC_GOOGLE_CLIENT_ID: z.string().min(1),
  },
  runtimeEnv: process.env,
  emptyStringAsUndefined: true,
})
