/**
 * ADB 操作相关
 * @returns 
 */

import { StartServer, StopServer, GetStatus } from "../../wailsjs/go/upload/Web"

export const useADB = () => {

    // 检测 ADB 是否存在
    const adbExists = () => {

    }
    return {
        adbExists
    }
}