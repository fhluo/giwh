export namespace main {
	
	export class Progress {
	    rarity: string;
	    count: number;
	
	    static createFrom(source: any = {}) {
	        return new Progress(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.rarity = source["rarity"];
	        this.count = source["count"];
	    }
	}
	export class StatResult {
	    wishType: string;
	    progresses: Progress[];
	
	    static createFrom(source: any = {}) {
	        return new StatResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.wishType = source["wishType"];
	        this.progresses = this.convertValues(source["progresses"], Progress);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

