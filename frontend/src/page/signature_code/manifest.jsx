import { observer, useStore } from "../../hooks/storeHook";

import NoApk from "../../component/NoApk";
import MenifestComp from "../../component/Manifest";


// 处理 apk 文件打包页面
const Manifest = () => {

    const { appStore } = useStore()
    return appStore.path ? <MenifestComp /> : <NoApk />
}

export default observer(Manifest)