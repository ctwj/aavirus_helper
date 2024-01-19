import { useStore, observer } from "../../hooks/storeHook";

import File from './file'
import Manifest from './manifest'
import Smail from "./smail";


// 特征码定位
const SignatureCode = () => {

    const { appStore } = useStore()

    if (appStore.func === 'file') { // 文件处理
        return <File /> 
    }

    if (appStore.func === 'manifest') { // 配置处理
        return <Manifest />
    }

    return <Smail /> // 代码处理

}

export default observer(SignatureCode)