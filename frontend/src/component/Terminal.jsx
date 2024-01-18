import React from "react";
import { observer } from "../hooks/storeHook";
import Terminal, { ColorMode } from 'react-terminal-ui';

import './Terminal.css'

import { useStore } from "../hooks/storeHook";

/**
 * Termial 显示命令日志， 
 * @param {*} props 
 * @returns 
 */
const TerminalUI = (props) => {
    const title = props?.title ?? 'Command Logs'
    const { appStore } = useStore()

    const handleTerminalInput = (terminalInput) => {
        if (props?.commandHandler) {
            props.commandHandler(terminalInput)
        }
    }
    return (
        <div style={{ height: '100%' }}>
            <Terminal style={{ height: '100%' }} name={title} colorMode={ ColorMode.Dark }  
                onInput={ handleTerminalInput }>
                { appStore.terminalLineData }
            </Terminal>
        </div>
    )
}

export default observer(TerminalUI)