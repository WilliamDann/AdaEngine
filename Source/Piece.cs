namespace Chess;

// piece color in chess
public enum Color
{
    White = 1,
    Black = 0
}

public enum PieceType
{
    None         = '_', // empty square    
    Invalid      = 'x', // square outside of the board

    BlackPawn    = 'p',
    BlackRook    = 'r',
    BlackKnight  = 'n',
    BlackBishop  = 'b',
    BlackQueen   = 'q',
    BlackKing    = 'k',

    WhitePawn    = 'P',
    WhiteRook    = 'R',
    WhiteKnight  = 'N',
    WhiteBishop  = 'B',
    WhiteQueen   = 'Q',
    WhiteKing    = 'K',
}

public class Piece
{
    public static Color GetPieceColor(PieceType piece)
    {
        if (Char.IsUpper((char)piece))
            return Color.White;
        return Color.Black;
    }
}