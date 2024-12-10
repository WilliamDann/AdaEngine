namespace Tests;

[TestClass]
public class TestChessboard
{
    [TestMethod]
    public void TestClearBoard()
    {
        Chessboard b = new();
        b.Clear();

        for (int x = -2; x < 10; x++)
            for (int y = -1; y < 10; y++)
                if (x < 0 || x >= 8 || y < 0 || y >= 8)
                    Assert.AreEqual(b.GetSquare(x, y).piece, PieceType.Invalid);
                else
                    Assert.AreEqual(b.GetSquare(x, y).piece, PieceType.None);
    }

    [TestMethod]
    public void TestBoardIter()
    {
        Chessboard b = new();

        foreach (Square square in b.SquareIter())
            Assert.AreNotEqual(square.piece, PieceType.Invalid);
    }

    [TestMethod]
    public void TestPlace()
    {
        Chessboard b = new();

        for (int x = 0; x < 8; x++)
            for (int y = 0; y < 8; y++)
                b.SetSquare(new Square(PieceType.Pawn, Color.White), x , y);

        foreach (Square sq in b.SquareIter())
            Assert.AreEqual(new Square(PieceType.Pawn, Color.White).Code(), sq.Code());
    }

    [TestMethod]
    public void TestClearSquare()
    {
        Chessboard b = new();

        for (int x = 0; x < 8; x++)
            for (int y = 0; y < 8; y++)
                b.SetSquare(new Square(PieceType.Pawn, Color.White), x , y);

        for (int x = 0; x < 8; x++)
            for (int y = 0; y < 8; y++)
                b.ClearSquare(x, y);

        for (int x = 0; x < 8; x++)
            for (int y = 0; y < 8; y++)
            {
                Color c = Color.Black;
                if (x % 2 == 0 ^ y % 2 == 0)
                    c = Color.White;

                Assert.AreEqual(new Square(PieceType.None, c).Char(), b.GetSquare(x, y).Char());
            }
    }
}