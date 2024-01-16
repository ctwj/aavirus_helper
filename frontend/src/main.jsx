import React from 'react'
import {createRoot} from 'react-dom/client'
import 'reset-css';
import './style.css'
import App from './App'

import StoreContext from './context/storeContext';
import store from './store'

const container = document.getElementById('root')
const root = createRoot(container)

root.render(
    <React.StrictMode>
        <StoreContext.Provider value={store}>
            <App/>
        </StoreContext.Provider> 
    </React.StrictMode>
)
