import './app.css'
import {api} from './models'
import {useEffect, useState} from 'react'
import {Get5Stars, GetUIDs, GetWishName} from '../wailsjs/go/main/App'
import Sidebar from './components/sidebar'
import {repository} from '../wailsjs/go/models'
import Character from './components/character'
import {Tabs, TabsList, TabsTrigger} from './components/ui/tabs'

function App() {
    let [uidList, setUIDList] = useState([] as number[])
    let [currentUID, setCurrentUID] = useState(0)
    let [items, setItems] = useState([] as repository.Item[])

    let currentWishType = api.SCharacterEventWish


    useEffect(() => {
        GetUIDs().then(result => {
            setUIDList(result)
            if (uidList.length != 0) {
                setCurrentUID(uidList[1])
            }
        })
    })

    useEffect(() => {
        Get5Stars(currentUID, currentWishType).then(result => {
            setItems(result)
        })
    }, [currentUID, currentWishType])

    return (
        <div className="flex flex-row bg-white p-4 w-screen h-screen">
            <Sidebar currentUID={currentUID}/>
            <div className="flex flex-col px-5 w-full h-full">
                <Tabs>
                    <TabsList>
                        {
                            api.SharedWishTypes.map(wishType => {
                                    let [wishName, setWishName] = useState('')
                                    useEffect(() => {
                                        GetWishName(wishType).then(result => {
                                            setWishName(result)
                                        })
                                    })
                                    return (
                                        // currentWishType = wishType
                                        <TabsTrigger value={wishName}>{wishName}</TabsTrigger>
                                    )
                                }
                            )
                        }
                    </TabsList>
                </Tabs>

                <div className=" py-12 px-12 tracking-wider space-y-3 w-full">
                    <div className="flex flex-row flex-wrap w-full gap-x-6 gap-y-6 transition duration-200">
                        {
                            items.map(item => <Character icon={item.icon} name={item.name} pulls={item.pulls}/>)
                        }
                    </div>
                </div>
            </div>
        </div>
    )
}

export default App
