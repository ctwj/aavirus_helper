import { makeObservable, observable, action, computed } from "mobx"
import { makePersistable } from 'mobx-persist-store'

export class Menu {
    id = Math.random()
    
    menu = 'home'

    constructor() {
        makeObservable(this, {
            menu: observable,
            setMenu: action,
        })
        makePersistable(
            this, 
            { 
                name: 'menuStore', 
                properties: [
                    'menu',
                ], 
                storage: window.localStorage
            }
        );

    }

    get isHydrated() {
        return isHydrated(this);
    }

    async getStoredData() {
        return getPersistedStore(this);
    }

    setMenu (val) {
        this.menu =  val
    }

}