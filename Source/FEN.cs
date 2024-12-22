namespace Chess;

// class for reading and handling FEN Notation
public class FEN
{
    // input info
    public string notation;

    public const string CLEAR  = "8/8/8/8/8/8/8/8 w - - 0 1";
    public const string INITAL = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1";


    public FEN(string notation)
    {
        this.notation = notation;
    }

    // get the side to move
    public Color activeColor()
    {
        if (notation.Split(' ')[1] == "w")
            return Color.White;
        return Color.Black;
    }

    // get a string representing the castling ability for both sides
    //  upper = white, lower = black
    //  q = queenside, k = kingside
    public string castlingAbility()
    {
        return notation.Split(' ')[2];
    }

    // en passant target square
    // null if no en enpassant exists
    public Point? enPassantSquare()
    {
        string square = notation.Split(' ')[3];
        if (square == "-")
            return null;
        return new Point(square);
    }

    // halfmoves since the last pawn move or capture
    //  used for the 50 move rule
    public int halfmoveClock()
    {
        return Convert.ToInt32(notation.Split(' ')[4]);
    }

    // number of full moves in the game
    public int fullmoveClock()
    {
        return Convert.ToInt32(notation.Split(' ')[5]);
    }

    public string fen()
    {
        return this.notation;
    }

    public override string ToString()
    {
        return fen();
    }
}