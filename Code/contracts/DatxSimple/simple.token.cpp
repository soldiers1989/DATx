#include <DatxioLib/DatxioLib.hpp>

class simpletoken : public datxio::contract {
   public:
      simpletoken( account_name self )
      :contract(self),_accounts( _self, _self){}

      void transfer( account_name from, account_name to, uint64_t quantity ) {
         require_auth( from );

         const auto& fromacnt = _accounts.get( from );
         datxio_assert( fromacnt.balance >= quantity, "overdrawn balance" );
         _accounts.modify( fromacnt, from, [&]( auto& a ){ a.balance -= quantity; } );

         add_balance( from, to, quantity );
      }

      void issue( account_name to, uint64_t quantity ) {
         require_auth( _self );
         add_balance( _self, to, quantity );
      }

   private:
      struct account {
         account_name owner;
         uint64_t     balance;

         uint64_t primary_key()const { return owner; }
      };

      datxio::multi_index<N(accounts), account> _accounts;

      void add_balance( account_name payer, account_name to, uint64_t q ) {
         auto toitr = _accounts.find( to );
         if( toitr == _accounts.end() ) {
           _accounts.emplace( payer, [&]( auto& a ) {
              a.owner = to;
              a.balance = q;
           });
         } else {
           _accounts.modify( toitr, 0, [&]( auto& a ) {
              a.balance += q;
              datxio_assert( a.balance >= q, "overflow detected" );
           });
         }
      }
};

DATXIO_ABI( simpletoken, (transfer)(issue) )
