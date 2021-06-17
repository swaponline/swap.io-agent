import * as dotenv from 'dotenv'
dotenv.config({
  path:
    process.env.NODE_ENV === 'production'
      ? './.env.production'
      : './.env.development',
})
import http from 'http'
import express from 'express'
import { Server } from 'socket.io'
import { generateAccessToken } from './auth/generateAccessToken'
import { verifyAccessToken } from './auth/verifyAccessToken'
const app = express()
const server = http.createServer(app)
const io = new Server(server)

app.get('/', (req, res) => {
  res.sendFile(__dirname + '/index.html')
})
app.get('/getToken', (req, res) => {
  res.send(
    generateAccessToken({
      id: 0,
    })
  )
})

io.on('connection', async (socket) => {
  try {
    await verifyAccessToken(socket.handshake.auth?.token)
    console.log('a user connected')
  } catch (e) {
    socket.disconnect()
    console.log('a user not connected')
  }
})

server.listen(process.env.PORT, () => {
  console.log(`listening on *:${process.env.PORT}`)
})
