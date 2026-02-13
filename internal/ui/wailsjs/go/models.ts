export namespace engine {
	
	export class Undo {
	    Captured: number;
	    CastleRights: number;
	    EnPassant: number;
	    HalfMove: number;
	
	    static createFrom(source: any = {}) {
	        return new Undo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Captured = source["Captured"];
	        this.CastleRights = source["CastleRights"];
	        this.EnPassant = source["EnPassant"];
	        this.HalfMove = source["HalfMove"];
	    }
	}

}

