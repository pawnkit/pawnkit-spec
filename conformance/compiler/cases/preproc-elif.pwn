#define ENABLED 1

#if ENABLED == 0
    #error wrong branch
#elif ENABLED == 1
    #define RESULT 1
#endif

main()
{
    return RESULT;
}
