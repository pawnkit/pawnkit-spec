main()
{
    new value = 1;
    {
        new value = 2;
        if (value != 2) {
            return 1;
        }
    }
    return value == 1 ? 0 : 1;
}
