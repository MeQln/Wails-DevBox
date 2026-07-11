import { ImageRead, ImageConvert, ImageCompress } from '../../wailsjs/go/main/App'

export type ImageInfo = {
  width: number
  height: number
  format: string
  size_bytes: number
  data_base64: string
}

export const imageApi = {
  read: (sourcePath: string) =>
    ImageRead(sourcePath) as Promise<ImageInfo>,

  convert: (sourcePath: string, targetFmt: string, outputPath: string) =>
    ImageConvert(sourcePath, targetFmt, outputPath) as Promise<ImageInfo>,

  compress: (sourcePath: string, quality: number, outputPath: string) =>
    ImageCompress(sourcePath, quality, outputPath) as Promise<ImageInfo>,
}