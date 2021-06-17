import jwt from 'jsonwebtoken'
import { ITokenData } from './types'

export const generateAccessToken = (payload: ITokenData): string =>
  jwt.sign(payload, process.env.TOKEN_SECRET, { expiresIn: '1800s' })
