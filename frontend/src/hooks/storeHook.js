import { useContext } from "react";
import StoreContext from "../context/storeContext";
import { observer } from "mobx-react"

function useStore() {
    const store = useContext(StoreContext);
    return store
}

export {
    observer,
    useStore,
}