import './app.css'
import {useEffect, useState} from 'react'
import {GetAssets, GetDefaultSharedWish, GetItems, GetLanguage, GetLocale} from '../wailsjs/go/main/App'
import Sidebar from './components/sidebar'
import {Tabs, TabsContent, TabsList, TabsTrigger} from './components/ui/tabs'
import Item from './components/item'
import {api} from '@/lib/models'
import {ScrollArea} from './components/ui/scroll_area'

function App() {
    let [uidList, setUIDList] = useState([] as string[])
    let [uid, setUID] = useState('')
    let [items, setItems] = useState([] as api.Item[])
    let [currentItems, setCurrentItems] = useState([] as api.Item[])

    let [wishes, setWishes] = useState([] as api.Wish[])
    let [sharedWishes, setSharedWishes] = useState([] as api.Wish[])
    let [currentWish, setCurrentWish] = useState('')

    let [locale, setLocale] = useState({} as api.Locale)
    let [locales, setLocales] = useState(new Map<string, api.Locale>())
    let [assets, setAssets] = useState({} as api.Assets)
    let [lang, setLang] = useState('')

    useEffect(() => {
        GetDefaultSharedWish().then(result => {
            setCurrentWish(result)
        })
        GetItems().then(result => {
            setItems(JSON.parse(result))
        })
        GetAssets().then(result => {
            setAssets(JSON.parse(result))
        })
        GetLanguage().then(result => {
            setLang(result)
        })
    }, [])

    const getLocale = async (lang: string): Promise<api.Locale> => {
        if (locales.has(lang)) {
            return locales.get(lang)!
        }
        let result: api.Locale = JSON.parse(await GetLocale(lang))
        setLocales(locales.set(lang, result))
        return result
    }


    useEffect(() => {
        getLocale(lang).then(result => {
            setLocale(result)
        })
    }, [lang])

    useEffect(() => {
        setUIDList([...new Set(items.map(item => item.uid))])
        setUID(uidList[0])

        let languages = new Set(items.map(item => item.lang))
        for (const lang of languages) {
            getLocale(lang).then()
        }

    }, [items])

    useEffect(() => {
        setSharedWishes(locale.sharedWishes)
        setWishes(locale.wishes)
    }, [locale])

    useEffect(() => {
        setCurrentItems(items.filter(item => {
            return item.uid == uid
        }).reverse())
    }, [items, uid])

    function getIcon(item: api.Item): string {
        const itemLocale = locales.get(item.lang)!
        if (item.name in itemLocale.charactersInverse) {
            return assets.characters[itemLocale.charactersInverse[item.name]]
        }
        if (item.name in itemLocale.weaponsInverse) {
            return assets.weapons[itemLocale.weaponsInverse[item.name]]
        }
        return ''
    }

    return (
        <div className="flex flex-row p-4 w-screen h-screen">
            <Sidebar uid={uid} setUID={setUID} uidList={uidList}/>
            <div className="flex flex-col px-5 w-4/5 h-full">
                <Tabs defaultValue={currentWish} onValueChange={value => {
                    let wishType = sharedWishes?.find(v => v.key.toString() == value)
                    setCurrentWish(wishType!.key)
                }} className="h-32 grow">
                    <TabsList className="bg-slate-200">
                        {
                            sharedWishes?.map(wishType => {
                                    return (
                                        <TabsTrigger value={wishType.key.toString()}
                                                     key={wishType.key}>{wishType.name}</TabsTrigger>
                                    )
                                }
                            )
                        }
                    </TabsList>
                    {
                        sharedWishes?.map(wishType => {
                            // console.log(items)
                            return (
                                <TabsContent value={wishType.key.toString()} className="bg-white/50" key={wishType.key}>
                                    <ScrollArea>
                                        <div
                                            className="flex flex-row flex-wrap w-full gap-x-6 gap-y-6 transition duration-200">
                                            {
                                                currentItems.filter(item => {
                                                    return item.rank_type == api.FiveStar && (
                                                        item.gacha_type == wishType.key || (
                                                            wishType.key == api.CharacterEventWish && item.gacha_type == api.CharacterEventWish2
                                                        )
                                                    )
                                                }).map(item =>
                                                    <Item item={item} pulls={0} key={item.id} icon={getIcon(item)}/>)
                                            }
                                        </div>
                                    </ScrollArea>
                                </TabsContent>
                            )
                        })
                    }
                </Tabs>
            </div>
        </div>
    )
}

export default App
