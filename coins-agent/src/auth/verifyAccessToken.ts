import jwt from 'jsonwebtoken'
import { ITokenData } from './types'

export const verifyAccessToken = (token: string): Promise<ITokenData> => {
  return new Promise((res, rej) => {
    jwt.verify(token, process.env.TOKEN_SECRET, (err, payload: ITokenData) => {
      if (err) return rej(new Error('incorrect token'))
      res(payload)
    })
  })
}
