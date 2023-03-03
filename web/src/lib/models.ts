export namespace api {
    export const NoviceWishes = '100'
    export const PermanentWish = '200'
    export const CharacterEventWish = '301'
    export const CharacterEventWishAndCharacterEventWish2 = '301'
    export const WeaponEventWish = '302'
    export const CharacterEventWish2 = '400'

    export const OneStar = '1'
    export const TwoStar = '2'
    export const ThreeStar = '3'
    export const FourStar = '4'
    export const FiveStar = '5'

    export interface Assets {
        characters: { [key: string]: string };
        weapons: { [key: string]: string };
    }

    export interface Item {
        count: string;
        gacha_type: string;
        id: string;
        item_id: string;
        item_type: string;
        lang: string;
        name: string;
        rank_type: string;
        time: string;
        uid: string;
    }

    export interface Language {
        key: string;
        name: string;
        short: string;
    }

    export interface Locale {
        language: Language;
        characters: { [key: string]: string };
        charactersInverse: { [key: string]: string };
        weapons: { [key: string]: string };
        weaponsInverse: { [key: string]: string };
        wishes: Wish[];
        sharedWishes: Wish[];
    }

    export interface Wish {
        key: string;
        name: string;
    }
}



