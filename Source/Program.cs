namespace Engine;
using Chess;

public class Program
{
    public static void Main()
    {
        FEN starting = new FEN("r1b2rk1/ppqn1ppp/2pb4/2p1p3/2N1P3/1P1P1N2/PBP2PPP/R2QK2R w KQ - 6 10");
        Chessboard b = starting.Board();

        Console.WriteLine(b);
        Console.WriteLine(starting);
        Console.WriteLine(starting.activeColor());
        Console.WriteLine(starting.castlingAbility());
        Console.WriteLine(starting.enPassantSquare());
        Console.WriteLine(starting.halfmoveClock());
        Console.WriteLine(starting.fullmoveClock());
    }
}