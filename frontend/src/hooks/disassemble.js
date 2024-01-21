import { Disassemble, FileList as AndroidManifestInfo, BatchPack, OpenOutput } from '../../wailsjs/go/project/Project'
import { useStore } from "../hooks/storeHook";
import { useState } from "react";

import { SplitButtonGroup, Dropdown, Spin, Toast, Divider, Button, Typography } from '@douyinfe/semi-ui';

export const useDisassemble = () => {

    const { appStore } = useStore()
    // 设置正在汇编的状态
    const [disassembling, setDisassembling] = useState()

    // 删除文件夹后打包
    const packApkWithDeleteDir = (mode) => {
        if (appStore.packing) {
            Toast.info({
                content: '请等待任务完成后再试',
                duration: 3,
            })
            return {};
        }

        appStore.setPacking(true)
        return BatchPack(appStore.disassembleDir, appStore.selFiles, mode).then(result => {
            appStore.setPacking(false)
            Toast.info({
                content: '打包完成',
                duration: 3,
            })
            OpenOutput()
        })
    }

    // 反汇编
    const  disassembleApk = () => {
        setDisassembling(true)  // 开启编译中状态
        return Disassemble(appStore.path).then((result) => { // 反编译完成后， 设置状态
            setDisassembling(false)
            appStore.setDisassembled(true)
            appStore.setDisassembleDir(result.outdir)
            return result.outdir
        });
    }

    return {
        packApkWithDeleteDir,
        disassembleApk,
        disassembling,
        setDisassembling,
    }
}