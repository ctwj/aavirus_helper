import { useStore, observer } from "../hooks/storeHook";

// 主面板显示的内容
const Manifest = () => {

    const { appStore } = useStore()

    return <div>
        Manifest
    </div>   
}

export default observer(Manifest)