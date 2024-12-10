using System.Security.Cryptography.X509Certificates;

namespace Chess;

// piece color in chess
public enum Color
{
    White = 1,
    Black = 0
}

public enum PieceType
{
    None    = '-', // empty square    
    Invalid = 'x', // square outside of the board

    Pawn    = 'p',
    Rook    = 'r',
    Knight  = 'n',
    Bishop  = 'b',
    Queen   = 'q',
    King    = 'k',
}

public class Square
{
    public PieceType piece;

    public Color color;

    public Square(PieceType piece, Color color)
    {
        this.piece = piece;
        this.color = color;
    }

    // get the code for the given piece
    public int Code()
    {
        return (int)this.color & (int)this.piece;
    }

    // get a character representation of the square
    public char Char()
    {
        char ch = (char)this.piece;

        if (this.color != Color.White)
        {
            if (this.piece == PieceType.None)
                ch = '+';
            else if (this.piece == PieceType.Invalid)
                ch = 'x';
            else
                ch -= (char)32; // make caps
        }

        return ch;
    }

    public override string ToString()
    {
        return ""+this.Char();
    }
}