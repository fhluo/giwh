interface AvatarProps {
    icon: string;
    name: string;
    pulls: number;
}

function Character({icon, name, pulls}: AvatarProps) {
    return (
        <div
            className="cursor-pointer flex flex-col w-fit items-center select-none shadow rounded-lg hover:bg-white/25 transition duration-200">
            <div>
                <div className="w-24">
                    <img alt={name} className="pointer-events-none rounded-t-lg shadow-inner bg-amber-600"
                         src={icon}/>
                </div>
            </div>
            {/*<div className="font-semibold text-gray-900 tracking-wide leading-relaxed">{name}</div>*/}
            <div
                className="text-sm text-gray-700 border-t py-0.5 w-full text-center rounded-b-lg leading-relaxed tracking-wider">{pulls}</div>
        </div>
    )
}

export default Character