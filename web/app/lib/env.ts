import { createEnv } from '@t3-oss/env-core'
import { z } from 'zod'

export const env = createEnv({
  clientPrefix: 'VITE_',
  server: {},
  client: {
    VITE_API_URL: z.string().url().default('http://localhost:8080'),
  },
  runtimeEnv: import.meta.env,
  skipValidation: import.meta.env.SKIP_ENV_VALIDATION === 'true',
  emptyStringAsUndefined: true,
})
